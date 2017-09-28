class Invoice

  QUEUE_NAME = "backend.invoices"

  def initialize(ch)
    @ch = ch
  end

  def start
    @q = @ch.queue(QUEUE_NAME)
    @x = @ch.default_exchange

    @q.subscribe do |delivery_info, properties, payload|
      r = begin
        data = JSON.parse(payload)

        invoices  = if data['client_id']
          self.class.invoices[data['client_id']] || []
        else
          self.class.invoices.values.flatten
        end
        { invoices: invoices }
      rescue JSON::ParserError
        { error: "Can't parse the payload - send {\"client_id\": \"<desired client id>\" } as payload" }
      end

      @x.publish(JSON.generate(r), :routing_key => properties.reply_to, :correlation_id => properties.correlation_id)

    end
    puts "Invoice provider started!"
  end

  def self.invoices
    {
      "google" => [
        {"total" => "200000000 USD", "services" => "Providing users data to the fbi", "customer" => "Federal Bureau of Investigation"},
        {"total" => "4000 USD", "services" => "Selling out users emails to ad companies", "customer" => "Big Bad Corporations"}
      ],
      "facebook" => [
        {"total" => "899229199 USD", "services" => "Selling out signed users photos", "customer" => "NSA"},
        {"total" => "9123 USD", "services" => "", "customer" => "Big Bad Corporation"}
      ],
      "yahoo" => [
        {"total" => "121223 USD", "services" => "Selling out the few emails they got", "customer" => "Big Bad Corporations"}
      ],
      "twitter" => [
        {"total" => "12 USD", "services" => "Telling everybody its hard to add edit button", "customer" => "Who cares"}
      ]
    }
  end

end
