FROM ubuntu:latest

# Arguments, default value ENV
ARG DBDIR=db

# Enviroments
ENV TODO_PORT=""
ENV TODO_DBFILE="db/scheduler.db"
ENV TODO_PASSWORD=""


# Workdir
WORKDIR /web
# Copy server and web dir
COPY todo_srv /web/
COPY web/ /web/web/
RUN mkdir /web/$DBDIR

CMD ["sh", "-c", "/web/todo_srv"]
