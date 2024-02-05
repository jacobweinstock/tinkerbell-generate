FROM scratch

COPY tinkerbell-generate /tinkerbell-generate

ENTRYPOINT [ "/tinkerbell-generate" ]