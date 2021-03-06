# SAML AND OIDC AUTHENTICATION PROTOTYPE

## INTRODUCTION
The aim of this prototype is to show how both SAML and OIDC authentication mechanisms can coexist in the same solution, whether SAML is used to authenticate the access to a web application (i.e. gosaml) while OIDC is used to authenticate an API call to another application (i.e. goapi).
Both authentication mechanisms are implemented using [Keycloak](https://www.keycloak.org/), but this solution includes more components such as [Nginx](https://www.nginx.com/) for path-based routing, [Keycloak Gatekeeper](https://github.com/keycloak/keycloak-gatekeeper) as OIDC reverse-proxy and [Kong](https://konghq.com/) as API Gateway.
To better explain the role that each component plains in this example solution, we'll adopt a step-by-step approach, by introducing each element one by one.  
To make it simple, the first stages don't implement HTTPS, that is introduced in stage **6-https**.
The following pictures shows these components and how they are connected together:

!["TECHNOLOGY VIEW"](docs/SingleSignOnExample-TechnologyView.png)

## DEPENDENCIES
This example is based on a single Docker Compose file that includes all the services involved at each stage of the implementation. 
Use `git checkout` to switch from one stage to another, accordingly to each corresponding tag:
- 1-samlonly
- 2-oidc
- 3-nginx
- 4-kong
- 5-kkfix
- 6-https

The following pictures shows how the container are instantiated in the final **docker-compose.yml** file and how their ports are mapped.
The browser interacts with ports 8000 and 8080, but the other ports can be used for test purposes.

!["TECHNOLOGY VIEW"](docs/SingleSignOnExample-Docker-Compose.png)


### **DISCLAIMER**
**Please be aware that this solution is implemented for test and education purposes only and it should not be uses ad-is in production environments. In particular, to simplify the setup, secrets and passwords are stored in the configuration files, without any specific protections that, instead, are a 'must-have' requirement for production environments.**

## SECURE CONNECTIONS
HTTPS has been introduced in stage **6-https** where keys and certificates are generated the first time the container **gjuljo/gencerts** is built.
The first time you run this container, typically at the first startup of the solution, you should find all the generated certificates in the `certs` directory.

Please trust the `myca.crt` at operating system level to see only **secure connections** in your browser (be aware that Firefox requires you to explicitly import the certificate).
On Windows, for instance, you can double-click the `myca.crt` file --> Install Certificate... --> Current User --> Place the certificate in the *Trusted Root Certificate Authorities* store.

## STARTUP
As Keycloak takes a while to startup, it's better to start it before the other containers and wait till it's up and running:

```console
docker-compose up -d test-oidc
```

Ensure that Keycloak is up and running by looking its logs:

```console
docker-compose logs -f test-oidc
```

You should see something like the following:

```
oidc | 08:29:17,535 INFO  [org.jboss.as] (Controller Boot Thread) WFLYSRV0060: Http management interface listening on http://127.0.0.1:9990/management
oidc | 08:29:17,535 INFO  [org.jboss.as] (Controller Boot Thread) WFLYSRV0051: Admin console listening on http://127.0.0.1:9990
oidc | 08:29:17,536 INFO  [org.jboss.as] (Controller Boot Thread) WFLYSRV0025: Keycloak 7.0.0 (WildFly Core 9.0.2.Final) started in 42509ms - Started 682 of 987 services (701 services are lazy, passive or on-demand)
```

You can also open the Keycloak Admin console at the link [http://oidc.lvh.me:8080/auth/](http://oidc.lvh.me:8080/auth/) or [https://oidc.lvh.me:8443/auth/](https://oidc.lvh.me:8443/auth/) (if Secure Connections are enabled).

Once **Keycloak** is up and running, you can start all the other containers:

```console
docker-compose up -d
```

After a while you can open the application page at the link [http://gosaml.lvh.me:8000/](http://gosaml.lvh.me:8080/) or [https://gosaml.lvh.me:8000/](https://gosaml.lvh.me:8000/) (if Secure Connections are enabled).

You'll be redirected to the **Keycloak** login page (see the available user credentials below). Once logged in, the actual application page is shown and you can enter any sentence in the shown text field: every time you press the `OK` button, the order if the words in the sentence is reversed (i.e. "hello world" becomes "world hello").

## USERS AND PASSWORDS
The following users are defined in the Keycloak server:
- admin/password
- pippo/pippopippo
- paperino/paperinopaperino
- testotp/testotptestotp

The user `testotp`, in addition to the password, requires a **Time-based One Time Password** (TOTP) 6-digits code, generated by an Authenticator App seeded with following QR Code:

!["OTP-QRCODE"](docs/qrcode-otp-demo.png)

This QR Code corresponds to the URI:

```
otpauth://totp/demo:testotp?secret=NFZTCZLMOA3TERTSPJ3G2S2ENE2WONLV&issuer=demo&algorithm=SHA1&digits=6&period=30
```

where:
- label = demo:testotp
- issuer = demo
- algorithm = SHA1
- secret = NFZTCZLMOA3TERTSPJ3G2S2ENE2WONLV
- period = 30s
- digits = 6


## 1-SAMLONLY
Start the SAML application (gosaml).

```
docker-compose up -d test-gosaml
```

Once the application is started, open the browser at the link [http://gosaml.lvh.me:8000/](http://gosaml.lvh.me:8000/). You'll be redirected to the Keycloak login page as the application requires a SAML-based authentication.
Provide username and password (pippo/pippopippo) to login and the browser should open a simple web page where you can enter some text. When you press the **OK** button, the web application sends the text to the server that, in this implementation, will just reply with the text *this is the default implementation*.
In further stages of this solution, this back-end API will be replaced with a different implementation that will revert the word order of the text provided in the input field. The Nginx component will route the API call to this second implementation.

!["1-SAMLONLY"](docs/SingleSignOnExample-Sequence-1.png)

## 2-OIDC
In this second stage we introduce **Keycloak Gatekeeper** as security proxy, to inforce the OIDC authentication before calling the SAML application. Behind the scene a double authentication in in place, involving a single sign-on scenario, even if the user won't notice it as only one login is requested.

Stop the further solution, if needed, then start this new one once test-oidc is running again:

```
docker-compose up -d
```

The following sequence diagram shows the interation between the three involved components:

!["2-OIDC"](docs/SingleSignOnExample-Sequence-2.png)


## 3-NGINX
In the third stage we introduct two additional components:
- **goapi** that implements the revert API by actually reverting the word order of the sentence writte in the input field of the webpage
- **nginx** that routes any call to the revert API to the goapi component instead of the gosaml

The API call will actually include the OIDC Bearer token but no actual authentication check is done in goapi. For this reason in the next stage we'll introduce the **Kong** API Gateway to enforce the authentication, assuming goapi is not directly protected by the gatek reverse proxy.

The following sequence diagram shows the interaction between all the involved components at this stage:

!["3-OIDC"](docs/SingleSignOnExample-Sequence-3.png)


## 4-KONG
This is stage where also **Kong** API Gateway is introduced to protect the access to the goapi application, by verifying the OIDC Bearer token that the gatek reverse proxy introduces in the **Authorization** HTTP header.
Since the basic **jwt** Kong plugin is used, we have to declare Keycloak's public key in the **kong.yml** file.

The following sequence diagram shows the interaction of all the components of this solution:

!["4-KONG"](docs/SingleSignOnExample-Sequence-4.png)

## 5-KKFIX
In this stage you can find a fix for **Keycloak Gatekeeper** to better handle expired refresh tokens.

## 6-HTTPS
In this stage both **Keycloak** and **Keycloak Gatekeeper** expose HTTPS ports and the communication between them and with **Nginx** are secure.

Please trust the `myca.crt` certificate as described in the [Secure Connections](#-SECURE-CONNECTIONS) section.

