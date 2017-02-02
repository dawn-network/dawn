FROM scratch
COPY ./glogchain .
EXPOSE 46658 46657 46656
VOLUME db.db
VOLUME state.db
CMD ['glogchain']
