plugins {
    alias(libs.plugins.kotlin.jvm)
    alias(libs.plugins.ktor)
}

group = "dev.praja.porgs"
version = "0.0.1"

application {
    mainClass = "io.ktor.server.netty.EngineMain"
}

kotlin {
    jvmToolchain(25)
}

ktor {
    docker {
        jreVersion.set(JavaVersion.VERSION_25)
    }
}

dependencies {
    implementation(libs.kotlinx.serialization.json)
    implementation(libs.ktor.server.config.yaml)
    implementation(libs.ktor.server.core)
    implementation(libs.ktor.server.netty)
    implementation(libs.ktor.server.thymeleaf)
    implementation(libs.logback.classic)
    testImplementation(libs.kotlin.test.junit)
    testImplementation(libs.ktor.server.test.host)

    constraints {
        implementation("io.netty:netty-codec-http2:4.2.12.Final") {
            because("fixes CVE-2026-33870 found in 4.2.9.Final")
        }
    }
}

tasks.withType<JavaExec> {
    jvmArgs("--enable-native-access=ALL-UNNAMED")
}
