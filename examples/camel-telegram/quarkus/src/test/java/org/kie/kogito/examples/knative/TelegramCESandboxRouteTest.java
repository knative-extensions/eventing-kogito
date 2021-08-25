package org.kie.kogito.examples.knative;

import io.quarkus.test.junit.QuarkusTest;
import org.junit.jupiter.api.Test;

import static io.restassured.RestAssured.given;
import static org.hamcrest.CoreMatchers.containsString;

@QuarkusTest
public class TelegramCESandboxRouteTest {

    @Test
    public void testHelloEndpoint() {
        given()
                .when()
                .header("ce-type", "example")
                .header("ce-id",  "12345")
                .header("ce-source", "/local/unit/test")
                .body("{ \"data\": \"somethingValueable\" }")
                .post("/camel/ce")
                .then()
                .statusCode(200)
                .body(containsString("**ID:** 12345"));
    }
}
