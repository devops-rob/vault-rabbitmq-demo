job "consul" {

  datacenters = ["dc1"]
  
  group "HashiTalks" {
    count = 1
    network {
      mode = "host"
    }

    task "consul" {
      driver = "raw_exec"

      config {
          command = "consul"
          args    = ["agent", "-dev"] 
        }
      artifact {
          source = "https://releases.hashicorp.com/consul/1.6.3/consul_1.6.3_darwin_amd64.zip"
      }

    }
  }
}