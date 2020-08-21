require 'logger'
require 'contentful/webhook/listener'

require_relative 'publish_controller'
require_relative 'unpublish_controller'

Contentful::Webhook::Listener::Server.start do |config|
  config[:logger] = Logger.new(STDOUT)
  config[:endpoints] = [
    {
      endpoint: '/publish',
      controller: PublishController,
      timeout: 1
    },
    {
      endpoint: '/unpublish',
      controller: UnpublishController,
      timeout: 1
    }
  ]
end.join
