package dev.praja.porgs

import dev.praja.porgs.models.Person
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.thymeleaf.*
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.int
import kotlinx.serialization.json.jsonObject
import kotlinx.serialization.json.jsonPrimitive

fun Route.sandboxController(appConfig: AppConfig) {
    get("/sandbox") {
        val person = Person.find(1)

        call.respondTemplate(
            "sandbox/index",
            mapOf(
                "appConfig" to appConfig,
                "person" to person,
            )
        )
    }

    post("/sandbox/counter/increment") {
        // Read the current count signal sent by Datastar in the request body
        val body = call.receiveText()
        val json = Json.parseToJsonElement(body).jsonObject
        val current = json["count"]?.jsonPrimitive?.int ?: 0
        val next = current + 1

        call.response.header(HttpHeaders.CacheControl, "no-cache")
        call.respondTextWriter(contentType = ContentType.Text.EventStream) {
            write("event: datastar-patch-signals\n")
            write("data: signals {count: $next}\n")
            write("\n")
            flush()
        }
    }
}
