workflow "New workflow" {
  on = "push"
  resolves = ["dockerhub push"]
}

action "GitHub Action for Docker" {
  uses = "actions/docker/cli@76ff57a"
  args = "build -t jpweber/servicemon-operator ."
}

action "Docker Registry" {
  uses = "actions/docker/login@76ff57a"
  needs = ["GitHub Action for Docker"]
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "dockerhub push" {
  uses = "actions/docker/cli@76ff57a"
  needs = ["Docker Registry"]
  args = "push jpweber/servicemon-operator "
}
