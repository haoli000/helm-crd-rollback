FROM maven:3.9.9-eclipse-temurin-17 AS build
WORKDIR /app
COPY . .
RUN mvn clean package -DskipTests

FROM eclipse-temurin:17-alpine
WORKDIR /app
COPY --from=build /app/target/sample-leader-election-*-jar-with-dependencies.jar /app/app.jar

ENV POD_NAME=default-operator
ENV POD_NAMESPACE=default

ENTRYPOINT ["java", "-cp", "app.jar", "io.javaoperatorsdk.operator.sample.LeaderElectionTestOperator"]