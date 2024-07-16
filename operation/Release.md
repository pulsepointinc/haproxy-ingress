## Prepare an image

1. Make a branch
2. Update docker base images from registry.pulsepoint.com/haproxy:v2.4_ppfp-dev to a released haproxy image version
3. `make docker-build`
4. `docker tag localhost/haproxy-ingress:latest registry.pulsepoint.com/haproxy-ingress:v0.14_ppfp_VERSION`
5. `docker push JUST_BUILD_IMAGE`
