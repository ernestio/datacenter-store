# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

# spec/app_spec.rb
require File.expand_path '../spec_helper.rb', __FILE__

describe 'datacenters_data_microservice' do
  describe 'a non authorized access' do
    describe 'to create a datacenter' do
      it 'should throw a 403' do
        post '/datacenters'
        expect(last_response.status).to be 403
      end
    end
    describe 'get datacenter list' do
      it 'should throw a 403' do
        get '/datacenters'
        expect(last_response.status).to be 403
      end
    end
    describe 'get a specific datacenter' do
      it 'should throw a 403' do
        get '/datacenters/foo'
        expect(last_response.status).to be 403
      end
    end
  end

  describe 'an authorized access as admin' do
    let!(:username)  { 'username' }
    let!(:password)  { 'password' }
    let!(:client_id) { 'client_id' }
    let!(:admin)     { true }

    before do
      DatacenterModel.dataset.destroy
      @token = SecureRandom.hex
      AuthCache.set @token, { user_id:   SecureRandom.uuid,
                              client_id: client_id,
                              user_name: username,
                              admin:     admin }.to_json
      AuthCache.expire @token, 3600
    end

    describe 'create datacenter' do
      let!(:name) { 'foo' }
      let!(:data) do
        { client_id: client_id,
          datacenter_name: name,
          datacenter_type: 'vcloud',
          datacenter_region: 'XXX',
          datacenter_username: username,
          datacenter_password: 'password',
          vcloud_url: 'http://my.url.com',
          external_network: 'ext' }.to_json
      end

      before do
        post '/datacenters',
             data,
             'HTTP_X_AUTH_TOKEN' => @token,
             'CONTENT_TYPE' => 'application/json'
      end
      it 'should response with a 200' do
        expect(last_response.status).to be 200
      end
      it 'should persist the datacenter' do
        datacenters = DatacenterModel.dataset.filter(datacenter_name: name)
        expect(datacenters.count).to be 1
      end
      it 'should response the datacenter object' do
        datacenter = JSON.parse(last_response.body, symbolize_names: true)
        expect(datacenter[:datacenter_name]).to eq name
        expect(datacenter[:datacenter_username]).to eq username
        expect(datacenter[:datacenter_password]).to eq 'password'
      end
    end

    describe 'list datacenters' do
      before do
        1.upto(10) do |i|
          datacenter = DatacenterModel.new
          datacenter.datacenter_name = "datacenter_#{i}"
          datacenter.datacenter_id = i
          datacenter.client_id = client_id
          datacenter.datacenter_username = 'x'
          datacenter.datacenter_password = 'x'
          datacenter.datacenter_type = 'x'
          datacenter.datacenter_region = 'x'
          datacenter.vcloud_url = 'x'
          datacenter.external_network = 'x'
          datacenter.save
        end
        get '/datacenters', {}, 'HTTP_X_AUTH_TOKEN' => @token
      end

      it 'should response with a 200 code' do
        expect(last_response.status).to be 200
      end
      it 'should return a list of existing datacenters' do
        expect(JSON.parse(last_response.body).length).to be(10)
      end
    end

    describe 'get datacenter details' do
      before do
        datacenter = DatacenterModel.new
        datacenter.datacenter_name = 'foo'
        datacenter.datacenter_id = 'CBF9BD7C-C2BC-4DA1-A799-55149E631B0F'
        datacenter.client_id = 'client_id'
        datacenter.datacenter_username = 'username'
        datacenter.datacenter_password = 'password'
        datacenter.datacenter_type = 'x'
        datacenter.datacenter_region = 'x'
        datacenter.vcloud_url = 'x'
        datacenter.external_network = 'x'
        datacenter.save
        get '/datacenters/CBF9BD7C-C2BC-4DA1-A799-55149E631B0F/',
            {},
            'HTTP_X_AUTH_TOKEN' => @token
      end
      it 'should response with a 200 code' do
        expect(last_response.status).to be 200
      end
      it 'should return datacenter details' do
        datacenter = JSON.parse(last_response.body)
        datacenter[:datacenter_name] = 'foo'
        datacenter[:datacenter_username] = 'username'
        datacenter[:datacenter_password] = 'password'
      end
    end
  end

  describe 'an authorized access' do
    let!(:username)  { 'username' }
    let!(:password)  { 'password' }
    let!(:admin)     { false }
    let!(:client_id) { 'client_id' }

    before do
      DatacenterModel.dataset.destroy
      @token = SecureRandom.hex
      AuthCache.set @token, { user_id:   SecureRandom.uuid,
                              client_id: client_id,
                              user_name: username,
                              admin:     admin }.to_json
      AuthCache.expire @token, 3600
    end

    describe 'create datacenter' do
      let!(:name) { 'foo' }
      let!(:data) do
        { client_id: client_id,
          datacenter_name: name,
          datacenter_type: 'vcloud',
          datacenter_region: 'XXX',
          datacenter_username: 'username',
          datacenter_password: 'password',
          vcloud_url: 'http://my.url.com',
          external_network: 'ext' }.to_json
      end

      before do
        post '/datacenters', data, 'HTTP_X_AUTH_TOKEN' => @token
      end
      it 'should response with a 200' do
        expect(last_response.status).to be 200
      end
      it 'should persist the datacenter' do
        datacenters = DatacenterModel.dataset.filter(datacenter_name: name)
        expect(datacenters.count).to be 1
      end
      it 'should response the datacenter object' do
        datacenter = JSON.parse(last_response.body, symbolize_names: true)
        expect(datacenter[:datacenter_name]).to eq name
      end
    end

    describe 'list datacenters' do
      before do
        1.upto(10) do |i|
          datacenter = DatacenterModel.new
          datacenter.datacenter_name = "datacenter_#{i}"
          datacenter.datacenter_id = i
          datacenter.client_id = client_id
          datacenter.datacenter_username = 'x'
          datacenter.datacenter_password = 'x'
          datacenter.datacenter_type = 'x'
          datacenter.datacenter_region = 'x'
          datacenter.vcloud_url = 'x'
          datacenter.external_network = 'x'
          datacenter.save
        end
        get '/datacenters', {}, 'HTTP_X_AUTH_TOKEN' => @token
      end

      it 'should response with a 200 code' do
        expect(last_response.status).to be 200
      end
      it 'should return a list of existing datacenters' do
        expect(JSON.parse(last_response.body).length).to be(10)
      end
    end

    describe 'get datacenter details' do
      before do
        datacenter = DatacenterModel.new
        datacenter.datacenter_name = 'foo'
        datacenter.datacenter_id = '547697FE-D9EC-4CBA-AE39-AFF0F2323432'
        datacenter.client_id = client_id
        datacenter.datacenter_username = 'x'
        datacenter.datacenter_password = 'x'
        datacenter.datacenter_type = 'x'
        datacenter.datacenter_region = 'x'
        datacenter.vcloud_url = 'x'
        datacenter.external_network = 'x'
        datacenter.save
        get '/datacenters/547697FE-D9EC-4CBA-AE39-AFF0F2323432/',
            {},
            'HTTP_X_AUTH_TOKEN' => @token
      end
      it 'should response with a 200 code' do
        expect(last_response.status).to be 200
      end
      it 'should return datacenter details' do
        datacenter = JSON.parse(last_response.body)
        datacenter[:datacenter_name] = 'foo'
      end
    end

    describe 'search datacenter by name' do
      let!(:name) { 'foo' }
      before do
        datacenter = DatacenterModel.new
        datacenter.datacenter_name = 'foo'
        datacenter.datacenter_id = 'C598D850-E5D4-4B60-AE51-C52D68051CBE'
        datacenter.client_id = client_id
        datacenter.datacenter_username = 'x'
        datacenter.datacenter_password = 'x'
        datacenter.datacenter_type = 'x'
        datacenter.datacenter_region = 'x'
        datacenter.vcloud_url = 'x'
        datacenter.external_network = 'x'
        datacenter.save
        get '/datacenters/search/', { name: name }, 'HTTP_X_AUTH_TOKEN' => @token
      end
      it 'should response with a 200 code' do
        expect(last_response.status).to be 200
      end
      it 'should return datacenter details' do
        datacenter = JSON.parse(last_response.body)
        datacenter[:datacenter_name] = name
        datacenter[:datacenter_id] = 'C598D850-E5D4-4B60-AE51-C52D68051CBE'
      end
      describe 'for an unexisting datacenter' do
        let!(:name) { 'foobar' }
        it 'should response with a 404 code' do
          expect(last_response.status).to be 404
        end
      end
    end

    describe 'update a datacenter' do
      before do
        put '/datacenters/foo', '', 'HTTP_X_AUTH_TOKEN' => @token
      end
      it 'should return a Not Implemented response' do
        expect(last_response.status).to be(405)
      end
    end

    describe 'update a datacenter' do
      before do
        delete '/datacenters/foo', '', 'HTTP_X_AUTH_TOKEN' => @token
      end
      it 'should return a Not Implemented response' do
        expect(last_response.status).to be(405)
      end
    end
  end
end
