package dev.praja.porgs

import io.ktor.server.application.*
import io.ktor.server.http.content.*
import io.ktor.server.routing.*
import io.ktor.server.thymeleaf.*

fun Application.configureRouting() {
    val appConfig = attributes[AppConfigKey]

    routing {
        staticResources("/assets", "assets")
        staticResources("/css", "css")
        staticResources("/javascript", "javascript")

        get("/") {
            call.respondTemplate("index", mapOf("appConfig" to appConfig))
        }
    }
}
