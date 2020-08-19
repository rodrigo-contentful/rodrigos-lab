require 'contentful/management'
require 'contentful/webhook/listener'

require_relative 'utils'

class PublishController < Contentful::Webhook::Listener::Controllers::WebhookAware
  attr_reader :client, :request, :response, :environment

  def initialize(*)
    super
    @client = Contentful::Management::Client.new(ENV['CF_TEST_CMA_TOKEN'], raise_errors: true)
  end

  def publish
    return unless authenticated?
    logger.debug "Successfully authenticated request"

    return unless webhook.entry? || webhook.asset?
    logger.debug "Webhook is for Entry or Asset"

    logger.info "Received #{webhook.entry? ? 'Entry' : 'Asset'} ID: '#{webhook.id}' for Space: '#{webhook.space_id}'"

    push_entry if webhook.entry?
    push_asset if webhook.asset?
  end

  protected

  def authenticated?
    !(request['Authorization'] =~ /Bearer #{ENV['AUTH_TOKEN']}/).nil?
  end

  def pre_perform(request, response)
    @request = request
    @response = response

    @environment = client.environments(request['X-Target-Space-Id']).find(request['X-Target-Environment-Id'])
    logger.debug "Loaded Environment ID '#{environment.id}' for Space ID '#{environment.space.id}'"

    super(request, response)
  rescue Contentful::Management::RateLimitExceeded => e
    sleep(e.reset_time) if e.reset_time?
    retry
  end

  private

  def push_entry
    upsert_entry_in_target_space! unless entry_overridden_in_target_space?
  end

  def entry_overridden_in_target_space?
    @target = client.entries(environment.space.id, environment.id).find(webhook.id)

    @target.fields.fetch(:overridden, false)
  rescue Contentful::Management::RateLimitExceeded => e
    sleep(e.reset_time) if e.reset_time?
    retry
  rescue Contentful::Management::NotFound
    logger.debug "Entry does not exist in target space."
    @target = nil
    false
  end

  def upsert_entry_in_target_space!
    created = false
    unless @target
      content_type = client.content_types(environment.space.id, environment.id).find(webhook.sys['contentType']['sys']['id'])
      @target = client.entries(environment.space.id, environment.id).create(content_type, id: webhook.id)
      created = true
    end

    @target.reload unless created

    webhook.fields.each do |field, locales|
      @target.send("#{Utils.snakify(field)}_with_locales=", locales)
    end

    sleep(1)

    @target.save
    logger.info "Updated '#{@target.id}' in Space '#{@target.space.id}'"
  rescue Contentful::Management::UnprocessableEntity => e
    logger.error e.response.load_json
  rescue Contentful::Management::BadRequest => e
    logger.error e.response.load_json
    retry
  rescue Contentful::Management::Conflict => e
    logger.error e.response.load_json
    retry
  rescue Contentful::Management::RateLimitExceeded => e
    sleep(e.reset_time.to_i) if e.reset_time?
    retry
  rescue HTTP::ConnectionError => e
    logger.error 'Something went wrong. Stacktrace: '
    logger.error e
    retry
  rescue Errno::ETIMEDOUT => e
    logger.error 'Something went wrong. Stacktrace: '
    logger.error e
    retry
  rescue Exception => e
    logger.error 'Something went wrong. Stacktrace: '
    logger.error e
    @response.body = "Something went wrong"
    @response.status = 500
  end

  def push_asset
    upsert_asset_in_target_space! unless asset_overridden_in_target_space?
  end

  def asset_overridden_in_target_space?
    @target = client.assets(environment.space.id, environment.id).find(webhook.id)

    false
  rescue Contentful::Management::NotFound
    logger.debug "Asset does not exist in target space."
    @target = nil
    false
  rescue Contentful::Management::RateLimitExceeded => e
    sleep(e.reset_time) if e.reset_time?
    retry
  end

  def files_from_localized_files(locales)
    locales.each_with_object({}) do |(locale, value), result|
      file = ::Contentful::Management::File.new
      file.properties[:contentType] = value['contentType']
      file.properties[:fileName] = value['fileName']
      file.properties[:upload] = "https:#{value['url']}"

      result[locale] = file
    end
  end

  def upsert_asset_in_target_space!
    created = false
    unless @target
      @target = client.assets(environment.space.id, environment.id).create(id: webhook.id)
      created = true
    end

    @target.reload unless created

    webhook.fields.each do |field, locales|
      locales = files_from_localized_files(locales) if field == 'file'

      @target.send("#{Utils.snakify(field)}_with_locales=", locales)
    end

    sleep(1)

    @target.save

    begin
      @target.process_file
    rescue Contentful::Management::BadRequest
      logger.debug "Asset already processed. Skipping."
    end

    logger.info "Updated '#{@target.id}' in Space '#{@target.space.id}'"
  rescue Contentful::Management::RateLimitExceeded => e
    sleep(e.reset_time.to_i) if e.reset_time?
    retry
  rescue Contentful::Management::UnprocessableEntity => e
    logger.error e.response.load_json
  rescue HTTP::ConnectionError => e
    logger.error 'Something went wrong. Stacktrace: '
    logger.error e
    retry
  rescue Errno::ETIMEDOUT => e
    logger.error 'Something went wrong. Stacktrace: '
    logger.error e
    retry
  rescue Exception => e
    logger.error 'Something went wrong. Stacktrace: '
    logger.error e
    @response.body = "Something went wrong"
    @response.status = 500
  end
end
