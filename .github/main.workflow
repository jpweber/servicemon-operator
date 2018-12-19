workflow "New workflow" {
  on = "push"
  resolves = ["dockerhub push"]
}

action "docker build" {
  uses = "actions/docker/cli@76ff57a"
  runs = "build -t jpweber/servicemon-operator ."
}

action "Docker Registry" {
  uses = "actions/docker/login@76ff57a"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
  needs = ["docker build"]
}

action "dockerhub push" {
  uses = "actions/docker/cli@76ff57a"
  needs = ["Docker Registry"]
  runs = "push jpweber/servicemon-operator "
}
