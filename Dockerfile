from ubuntu:16.04

RUN apt-get update
RUN apt-get install -y python

# Need to place this init-proc somewhere inside the container
COPY init-proc /srv/

# Run the init-proc as PID 1 with the actual process we want to run as an
# argument to it.
CMD ["./srv/init-proc", "sleep", "100"]
