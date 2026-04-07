#!/usr/bin/env kotlin
import java.io.File
import java.time.LocalDateTime
import java.time.format.DateTimeFormatter

// ## Output Markdown file
val now = LocalDateTime.now()
val timestampSuffix = now.format(DateTimeFormatter.ofPattern("yyyyMMddHHmmss"))
val fileName = "dump_src-$timestampSuffix.md"
val tmpDir = File("tmp")
if (!tmpDir.exists()) tmpDir.mkdirs()
val outputFile = File(tmpDir, fileName)

// ## Configuration
val projectRoot = File(".").canonicalFile
val sourceRoots = listOf("db", "gradle", "src")
val includeExtensions = setOf(
    "kt", "kts", "yaml", "yml", "html", "kte", "css",
    "js", "json", "toml", "properties", "xml", "sql", "jte"
)
val excludeDirs =
    setOf(".gradle", ".idea", ".git", ".kotlin", "build", "docs", "jte-classes", "node_modules", "scripts", "tmp", "lib")

// ## Document sections
val sections = listOf(
    Section("Build Configuration") { path ->
        path.endsWith(".gradle.kts") ||
                path.endsWith(".properties") ||
                path.endsWith("libs.versions.toml")
    },
    Section("Application Bootstrap") { path ->
        path.endsWith("Application.kt")
    },
    Section("Database Connection Setup") { path ->
        path.endsWith("Database.kt")
    },
    Section("Database Schema Migrations") { path ->
        path.startsWith("db/migrate")
    },
    Section("Exposed Tables") { path ->
        path.startsWith("src/main/kotlin/tables")
    },
    Section("Model Classes") { path ->
        path.startsWith("src/main/kotlin/models")
    },
    Section("Routing") { path ->
        path.endsWith("Routing.kt")
    },
    Section("Controllers") { path ->
        path.startsWith("src/main/kotlin/controllers")
    },
    Section("Stylesheets") { path ->
        path.startsWith("src/main/resources/css")
    },
    Section("JavaScript Components") { path ->
        path.startsWith("src/main/resources/javascript/components")
    },
    Section("JavaScript Code") { path ->
        path.startsWith("src/main/resources/javascript")
    },
    Section("Templates") { path ->
        path.startsWith("src/main/resources/templates")
    },
    Section("Other Resources") { path ->
        path.contains("src/main/resources")
    },
    Section("Kotlin Source Code") { path ->
        path.startsWith("src/main/kotlin") || path.startsWith("src/main/java")
    },
    Section("Tests") { path ->
        path.startsWith("src/test")
    },
    Section("Other") { _ -> true }
)

// ## Execute
val files = collectFiles()
val grouped = groupFiles(files)
outputFile.writeText(render(grouped))
println("✓ Written ${files.size} files to ${outputFile.path}")

// ## Helpers ——————————————————————————————
data class Section(val title: String, val match: (String) -> Boolean)
data class SourceFile(val relativePath: String, val content: String)

/** Generates a Markdown-compatible anchor link.
 * Example: "Build & Project Config" -> "build--project-config"
 */
fun String.toAnchor(): String =
    this.lowercase()
        .replace(Regex("[^a-z0-9 ]"), "") // Remove special chars
        .replace(" ", "-")                // Spaces to hyphens

fun languageFor(path: String): String = when (path.substringAfterLast('.', "")) {
    "kt", "kts" -> "kotlin"
    "xml" -> "xml"
    "sql" -> "sql"
    "yml", "yaml" -> "yaml"
    "jte" -> "html"
    else -> path.substringAfterLast('.', "")
}

// ## Logic ——————————————————————————————
fun collectFiles(): List<SourceFile> {
    val results = mutableMapOf<String, SourceFile>()

    // Files in the project root
    projectRoot.listFiles()?.filter {
        it.isFile && it.extension in includeExtensions
    }?.forEach { file ->
        results[file.name] = SourceFile(file.name, file.readText())
    }

    // Files in the source directories
    sourceRoots.map { projectRoot.resolve(it) }.filter { it.exists() }.forEach { root ->
        root.walkTopDown()
            .onEnter { it.name !in excludeDirs }
            .filter { it.isFile && it.extension in includeExtensions }
            .forEach { file ->
                val rel = file.relativeTo(projectRoot).path
                results[rel] = SourceFile(rel, file.readText())
            }
    }

    return results.values.sortedBy { it.relativePath }
}

fun groupFiles(files: List<SourceFile>): List<Pair<Section, List<SourceFile>>> {
    val buckets = sections.associateWith { mutableListOf<SourceFile>() }
    files.forEach { file ->
        val section = sections.first { it.match(file.relativePath) }
        buckets[section]!!.add(file)
    }
    return sections.map { it to buckets[it]!! }.filter { it.second.isNotEmpty() }
}

fun render(grouped: List<Pair<Section, List<SourceFile>>>): String = buildString {
    val displayTimestamp = now.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"))
    appendLine("# ${projectRoot.name} — Source Dump\n> Generated: $displayTimestamp\n")

    appendLine("## Table of Contents\n")
    grouped.forEach { (section, files) ->
        appendLine("- [${section.title}](#${section.title.toAnchor()})")
        files.forEach { file ->
            // Sub-links to individual files
            appendLine("  - [${file.relativePath}](#${file.relativePath.toAnchor()})")
        }
    }

    grouped.forEach { (section, files) ->
        appendLine("\n---\n## ${section.title}\n")
        files.forEach { file ->
            appendLine("### ${file.relativePath}\n")
            appendLine("```${languageFor(file.relativePath)}\n${file.content.trimEnd()}\n```\n")
        }
    }
}


