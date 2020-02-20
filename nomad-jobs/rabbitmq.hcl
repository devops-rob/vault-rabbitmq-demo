job "RabbitMQ" {

  datacenters = ["dc1"]
  type = "service"

  group "HashiTalks" {
    count = 1
    network {
      mode = "host"
    }

    update {
      max_parallel = 1
    }

    migrate {
      max_parallel = 1
      health_check = "checks"
      min_healthy_time = "5s"
      healthy_deadline = "30s"
    }

    task "rabbit" {
      driver = "docker"

      config {
        image = "rabbitmq:3-management"
        hostname = "${attr.unique.hostname}"
        port_map {
          http = 15672
          amqp = 5672
          ui = 15672
          epmd = 4369
          clustering = 25672
        }
      }
      resources {
        network {
          port "amqp" { static = 5672 }
          port "ui" { static = 15672 }
          port "epmd" { static = 4369 }
          port "clustering" { static = 25672 }
        }
      }

      service {
        name = "rabbitmq"
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