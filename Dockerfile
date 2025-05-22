FROM oven/bun:1.1-alpine

MAINTAINER Sasaya <sasaya@percussion.life>

COPY ./src ./src
COPY ./package.json ./package.json
COPY ./tsconfig.json ./tsconfig.json
COPY ./bun.lockb ./bun.lockb

RUN bun install

CMD ["bun", "run", "src/index.ts"]
