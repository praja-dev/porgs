package dev.praja

import io.ktor.server.application.*
import io.ktor.server.http.content.*
import io.ktor.server.routing.*
import io.ktor.server.thymeleaf.*

fun Application.configureRouting() {
    routing {
        staticResources("/css", "css")
        staticResources("/javascript", "javascript")

        get("/") {
            call.respondTemplate("index", mapOf("message" to "porgs: 0.0.1"))
        }
    }
}
