module Utils
  def self.snakify(object)
    snake = String(object).gsub(/(?<!^)[A-Z]/) { "_#{$&}" }
    snake.downcase
  end
end
