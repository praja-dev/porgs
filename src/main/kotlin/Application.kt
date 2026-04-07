package dev.praja.porgs

import io.ktor.server.application.*
import io.ktor.server.thymeleaf.*
import io.ktor.util.*
import org.thymeleaf.templateresolver.ClassLoaderTemplateResolver

data class AppConfig(
    val name: String,
    val shortName: String,
    val version: String,
)

val AppConfigKey = AttributeKey<AppConfig>("AppConfig")

fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}

fun Application.bootstrap() {
    val appConfig = AppConfig(
        name = environment.config.property("app.name").getString(),
        shortName = environment.config.property("app.shortName").getString(),
        version = environment.config.property("app.version").getString()
    )
    attributes.put(AppConfigKey, appConfig)

    val dbFile = "storage/porgs.db"
    val readPoolSize = environment.config.property("app.db.readPoolSize").getString().toInt()
    val busyTimeout = environment.config.property("app.db.busyTimeout").getString().toInt()

    java.io.File(dbFile).parentFile?.mkdirs()
    AppDatabase.init(dbFile, readPoolSize, busyTimeout)

    install(Thymeleaf) {
        setTemplateResolver(ClassLoaderTemplateResolver().apply {
            prefix = "templates/"
            suffix = ".html"
            characterEncoding = "utf-8"
        })
    }

    configureRouting()
}
