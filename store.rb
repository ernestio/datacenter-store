# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

require 'sequel'
require 'sinatra'
require 'flowauth'
require 'yaml'
require 'nats/client'

module Config
  def self.load_db
    NATS.start(servers: [ENV['NATS_URI']]) do
      NATS.request('config.get.postgres') do |r|
        return JSON.parse(r, symbolize_names: true)
      end
    end
  end
  def self.load_redis
    return if ENV['RACK_ENV'] == 'test'
    NATS.start(servers: [ENV['NATS_URI']]) do
      NATS.request('config.get.redis') do |r|
        return r
      end
    end
  end
end

class API < Sinatra::Base
  configure do
    # Default DB Name
    ENV['DB_URI'] ||= Config.load_db[:url]
    ENV['DB_REDIS'] ||= Config.load_redis
    ENV['DB_NAME'] ||= 'datacenters'

    #  Initialize database
    Sequel::Model.plugin(:schema)
    DB = Sequel.connect("#{ENV['DB_URI']}/#{ENV['DB_NAME']}")

    # Create datacenters database table if does not exist
    DB.create_table? :datacenters do
      String :datacenter_id, null: false, primary_key: true
      String :client_id, null: false
      String :datacenter_name, null: false
      String :datacenter_type, null: false
      String :datacenter_region, null: false
      String :datacenter_username, null: false
      String :datacenter_password, null: false
      String :vcloud_url, null: false
      String :vse_url, null: true
      String :external_network, null: false
      unique [:client_id, :datacenter_name, :datacenter_type]
    end
    # Create models and assign them to database tables
    Object.const_set('DatacenterModel', Class.new(Sequel::Model(:datacenters)))
  end

  before do
    content_type :json
  end

  use Authentication

  post '/datacenters/?' do
    datacenter = JSON.parse(request.body.read, symbolize_names: true)
    halt 400, 'Invalid external network' unless datacenter[:external_network]
    halt 400, 'Not vcloud url specified' unless datacenter[:vcloud_url]
    existing_datacenter = DatacenterModel.filter(datacenter_name: datacenter[:datacenter_name], client_id: env[:current_user][:client_id]).first
    unless existing_datacenter.nil?
      halt 409, url("/datacenters/#{existing_datacenter[:datacenter_id]}")
    end
    datacenter[:client_id] = env[:current_user][:client_id]
    datacenter[:datacenter_id] = SecureRandom.uuid
    DatacenterModel.insert(datacenter)
    datacenter.to_json
  end

  get '/datacenters/?' do
    fields = [:datacenter_id,
              :datacenter_name,
              :datacenter_type,
              :datacenter_region,
              :datacenter_username,
              :datacenter_password,
              :vcloud_url,
              :vse_url,
              :external_network]
    filters = { client_id: env[:current_user][:client_id] }
    DatacenterModel.select(*fields).filter(filters).all.map(&:to_hash).to_json
  end

  get '/datacenters/search/?' do
    datacenter = DatacenterModel.filter(client_id: env[:current_user][:client_id], datacenter_name: params[:name]).first
    halt 404 if datacenter.nil?
    status 200
    return { datacenter_id: datacenter[:datacenter_id],
             datacenter_name: datacenter[:datacenter_name],
             datacenter_type: datacenter[:datacenter_type],
             datacenter_region: datacenter[:datacenter_region],
             datacenter_username: datacenter[:datacenter_username],
             datacenter_password: datacenter[:datacenter_password],
             external_network: datacenter[:external_network],
             vcloud_url: datacenter[:vcloud_url],
             vse_url: datacenter[:vse_url] }.to_json
  end

  get '/datacenters/:datacenter/?' do
    fields = [:datacenter_id, :datacenter_name, :datacenter_type, :datacenter_region, :external_network, :vcloud_url]
    filters = { client_id: env[:current_user][:client_id], datacenter_id: params[:datacenter] }
    DatacenterModel.select(*fields).filter(filters).first.to_hash.to_json
  end

  put '/datacenters/:datacenter/?' do
    halt 405
  end

  delete '/datacenters/:datacenter/?' do
    halt 405
  end
end
