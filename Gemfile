source 'https://rubygems.org'

gem 'flowauth', path: '/opt/ernest-libraries/authentication-middleware'

gem 'sinatra'
gem 'sequel'
gem 'pg'
gem 'nats', git: 'https://github.com/r3labs/ruby-nats.git'

group :development, :test do
  gem 'pry'
  gem 'rack-test'
end

group :test do
  gem 'rspec'
  gem 'rubocop',   require: false
  gem 'simplecov', require: false
end
