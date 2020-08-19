require 'contentful/webhook/listener'

class UnpublishController < Contentful::Webhook::Listener::Controllers::WebhookAware
  def unpublish
    if webhook.entry?
      logger.info "published Entry ID: #{webhook.id} for Space: #{webhook.space_id}"
    end
  end

  private

  def unpublish_entry(entry)
  end

  def unpublish_asset(asset)
  end
end
