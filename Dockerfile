FROM library/alpine:3.1

RUN apk add --update \
	bash \
	curl \
	jq \
	git \
	tar

ADD ./backup-repos.sh /bin/backup-repos.sh
RUN chmod +x /bin/backup-repos.sh

ENTRYPOINT ["/bin/backup-repos.sh"]
