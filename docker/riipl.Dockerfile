# Build from the RIIPL library @ https://hub.docker.com/r/riipl/3d_qifp/
FROM riipl/3d_qifp

LABEL maintainer "Daniel Blezek blezek.daniel@mayo.edu"

# install grunt
WORKDIR /riipl
RUN mkdir -p /riipl/grunt_work
COPY bin/grunt-docker /bin/grunt
COPY docker/riipl.gruntfile.yml /riipl/gruntfile.yml
COPY docker/riipl.sh /riipl/riipl.sh

# What do we run on startup?
ENTRYPOINT ["/bin/grunt", "/riipl/gruntfile.yml"]
# We expose port 9901 by default
EXPOSE 9901:9901

