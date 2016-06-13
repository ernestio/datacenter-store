# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

require 'pry'
require 'simplecov'
require 'rack/test'

if ENV['COVERAGE']
  SimpleCov.start do
    minimum_coverage ENV['MIN_COVERAGE']
    add_filter '/spec/'
  end
  SimpleCov.at_exit do
    SimpleCov.result.format!
    if SimpleCov.result.covered_percent < SimpleCov.minimum_coverage
      covered  = format('%.2f', SimpleCov.result.covered_percent)
      min      = format('%.2f', SimpleCov.minimum_coverage)

      puts <<-eos
Coverage (#{covered}%) does not accomplish minimum #{min}%
----------------------------------------------------------

Please run on you dev environment:

:> COVERAGE=true bundle exec rspec -c -f d spec
:> open coverage/index.html

To see your coverage report.
Happy hacking
      eos

      exit(1)
    end
  end
end

ENV['DB_URI'] = 'postgres://postgres@127.0.0.1'
ENV['DB_NAME'] = 'test_datacenters'
ENV['NATS_URI'] = 'nats://127.0.0.1:4222'
ENV['RACK_ENV'] = 'test'

require_relative '../store'

module RSpecMixin
  include Rack::Test::Methods
  def app
    ::API
  end
end

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.include RSpecMixin
end
