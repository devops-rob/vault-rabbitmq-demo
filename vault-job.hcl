job "Vault" {

  datacenters = ["dc1"]
  type = "service"

  group "Vault" {
    count = 1

    update {
      max_parallel = 1
    }

    migrate {
      max_parallel = 1
      health_check = "checks"
      min_healthy_time = "5s"
      healthy_deadline = "30s"
    }

    task "vault" {
        driver = "docker"

        config {
            image = "vault:latest"
            hostname = "${attr.unique.hostname}"
            ipc_mode = "host"
        }
        env {
            VAULT_DEV_ROOT_TOKEN_ID  = "devopsrob"
            VAULT_DEV_LISTEN_ADDRESS = "0.0.0.0:8200"
        }
      resources {
        network {
          port "ui" { static = 8200 }
          port "clustering" { static = 8201 }
        }
      }

      service {
        name = "vault"
        port = "ui"
        check {
          name     = "alive"
          type     = "tcp"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}