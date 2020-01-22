# Use the standard Nginx image from Docker Hub
FROM openjdk:8-jre-slim AS sslbuilder
RUN apt-get update

WORKDIR /work/
VOLUME /ssl

COPY *.conf ./
COPY *.ext ./

RUN openssl req -x509 -subj "/CN=ca.lvh.me" -newkey rsa:2048 -nodes -keyout myca.key -out myca.crt                                                                                      \
 && keytool -genkeypair -alias server -keyalg RSA -keystore application.keystore -dname CN=oidc.lvh.me -ext san=dns:oidc.lvh.me -storepass password -keypass password -storetype pkcs12 \
 && keytool -certreq -alias server -keystore application.keystore -storepass password -ext san=dns:oidc.lvh.me > oidc.csr                                                               \
 && openssl x509 -req -in oidc.csr -CA myca.crt -CAkey myca.key -CAcreateserial -out oidc.crt -days 2048 -sha256 -extfile oidc.ext                                                      \          
 && keytool -import -trustcacerts -alias myca -keystore application.keystore -storepass password -file myca.crt -noprompt                                                               \
 && keytool -importcert -alias server -keystore application.keystore -storepass password -file oidc.crt                                                                                 \
 && openssl genrsa -out gatek.key 2048                                                                                                                                                  \
 && openssl req -new -sha256 -key gatek.key -subj "/CN=gosaml.lvh.me" -out gatek.csr -config gatek.conf                                                                                 \
 && openssl x509 -req -in gatek.csr -CA myca.crt -CAkey myca.key -CAcreateserial -out gatek.crt -days 2048 -sha256 -extfile gatek.ext                                                   \
 && openssl genrsa -out nginx.key 2048                                                                                                                                                  \
 && openssl req -new -sha256 -key nginx.key -subj "/CN=nginx" -out nginx.csr                                                                                                              \
 && openssl x509 -req -in nginx.csr -CA myca.crt -CAkey myca.key -CAcreateserial -out nginx.crt -days 2048 -sha256

CMD [ "cp", "-r", "/work/.", "/ssl/" ]