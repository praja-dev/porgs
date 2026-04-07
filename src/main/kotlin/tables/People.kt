package dev.praja.porgs.tables

import org.jetbrains.exposed.v1.core.dao.id.LongIdTable

object People : LongIdTable("people"){
    val name = varchar("name", 255)
}
