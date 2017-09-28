class Client

  QUEUE_NAME = "backend.clients"

  def initialize(ch)
    @ch = ch
  end

  def start
    @q = @ch.queue(QUEUE_NAME)
    @x = @ch.default_exchange

    @q.subscribe do |delivery_info, properties, payload|
      r = { clients: self.class.clients }
      @x.publish(JSON.generate(r), :routing_key => properties.reply_to, :correlation_id => properties.correlation_id)
    end
    puts "Client provider started!"
  end

  def self.clients
    ["google", "facebook", "yahoo", "twitter"]
  end

end
