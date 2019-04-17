FROM alpine:3.6
COPY ./bin/server /home/server/server

RUN apk --update add ca-certificates && \
		adduser -D server && \
		chown server /home/server/server && \
		chmod o+x /home/server/server

USER server
ENTRYPOINT /home/server/server