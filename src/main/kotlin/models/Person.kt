package dev.praja.porgs.models

import dev.praja.porgs.dbRead
import dev.praja.porgs.tables.People
import org.jetbrains.exposed.v1.core.eq
import org.jetbrains.exposed.v1.jdbc.selectAll

data class Person(
    val id: Long,
    val name: String,
) {
    companion object {
        suspend fun find(id: Long): Person? = dbRead {
            People.selectAll()
                .where { People.id eq id }
                .map { Person(it[People.id].value, it[People.name]) }
                .singleOrNull()
        }

        suspend fun all(): List<Person> = dbRead {
            People.selectAll()
                .map { Person(it[People.id].value, it[People.name]) }
        }
    }
}
