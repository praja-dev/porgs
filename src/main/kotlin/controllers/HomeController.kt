package dev.praja.porgs

import io.ktor.server.routing.*
import io.ktor.server.thymeleaf.*

fun Route.homeController(appConfig: AppConfig) {
    get("/") {
        call.respondTemplate(
            "home/index",
            mapOf(
                "appConfig" to appConfig,
            )
        )
    }
}
