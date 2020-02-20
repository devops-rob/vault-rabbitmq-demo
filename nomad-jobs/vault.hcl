job "Vault" {

  datacenters = ["dc1"]
  
  group "HashiTalks" {
    count = 1
    network {
      mode = "host"
    }

    task "vault" {
      driver = "raw_exec"

      config {
          command = "vault"
          args    = ["server", "-dev"] 
        }
      artifact {
          source = "https://releases.hashicorp.com/vault/1.3.2/vault_1.3.2_darwin_amd64.zip"
      }
      env {
          VAULT_DEV_ROOT_TOKEN_ID  = "devopsrob"
          VAULT_DEV_LISTEN_ADDRESS = "0.0.0.0:8200"
      }


    }
  }
}