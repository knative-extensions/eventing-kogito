package org.kie.kogito.examples.knative;

import javax.enterprise.context.ApplicationScoped;

import org.apache.camel.builder.endpoint.EndpointRouteBuilder;
import org.apache.camel.component.telegram.TelegramConstants;
import org.apache.camel.component.telegram.TelegramParseMode;

/**
 * Class to run locally the integration and check if everything is fine before deploying on Kubernetes
 * See ./hack/deploy-k8s.sh for how to deploy the actual route
 */
@ApplicationScoped
public class TelegramCESandboxRoute extends EndpointRouteBuilder {

    @Override
    public void configure() throws Exception {
        from(platformHttp("/camel/ce"))
                .convertBodyTo(String.class)
                .choice()
                    .when(header("ce-chat-id").isNull())
                    .setHeader("ce-chat-id", constant("{{defaultChatId}}"))
                .end()
                .setHeader(TelegramConstants.TELEGRAM_PARSE_MODE, constant(TelegramParseMode.MARKDOWN))
                .setHeader(TelegramConstants.TELEGRAM_CHAT_ID, header("ce-chat-id"))
                .to("string-template:templates/TelegramMessage.tm")
                // uncomment here to use send the message :)
                //.to("telegram:bots?authorizationToken={{authorizationToken}}")
                .marshal().json()
                // comment here to not log anymore
                .to(log("info").showBodyType(false).showBody(true).showExchangePattern(false));
    }
}
