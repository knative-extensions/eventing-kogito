import org.apache.camel.builder.RouteBuilder;
import org.apache.camel.component.telegram.TelegramConstants;
import org.apache.camel.component.telegram.TelegramParseMode;

public class TelegramCloudEventNotification extends RouteBuilder {

    @Override
    public void configure() throws Exception {
        from("knative:channel/kogito-channel")
                .convertBodyTo(String.class)
                .choice()
                    .when(header("ce-chat-id").isNull())
                    .setHeader("ce-chat-id", constant("{{defaultChatId}}"))
                .end()
                .setHeader(TelegramConstants.TELEGRAM_PARSE_MODE, constant(TelegramParseMode.MARKDOWN))
                .setHeader(TelegramConstants.TELEGRAM_CHAT_ID, header("ce-chat-id"))
                .to("string-template:file:/etc/camel/resources/TelegramMessage.tm")
                .to("telegram:bots?authorizationToken={{authorizationToken}}")
                .marshal().json()
                // comment here to not log anymore
                .to("log:info");
    }
}
