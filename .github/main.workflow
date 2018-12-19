workflow "New workflow" {
  on = "push"
  resolves = [
    "docker login",
    "GitHub Action for Docker",
  ]
}

action "docker login" {
  uses = "actions/docker/login@76ff57a"
  secrets = ["DOCKER_PASSWORD", "DOCKER_USERNAME"]
}

action "docker build1" {
  uses = "actions/docker/cli@76ff57a"
  args = ["build", "-t", "jpweber/servicemon-operator", "."]
}

action "GitHub Action for Docker" {
  uses = "actions/docker/cli@76ff57a"
  needs = ["docker build1", "docker login"]
  args = ["push", "jpweber/servicemon-operator"]
  secrets = ["DOCKER_PASSWORD", "DOCKER_USERNAME"]
}
