package dev.praja.porgs

import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import org.jetbrains.exposed.v1.core.Transaction
import org.jetbrains.exposed.v1.jdbc.Database
import org.jetbrains.exposed.v1.jdbc.transactions.transaction


object AppDatabase {
    lateinit var reader: Database private set
    lateinit var writer: Database private set

    val readDispatcher = Dispatchers.IO.limitedParallelism(10)
    val writeDispatcher = Dispatchers.IO.limitedParallelism(1)

    fun init(dbFile: String, readPoolSize: Int, busyTimeout: Int) {
        fun hikari(readOnly: Boolean, poolSize: Int) = HikariDataSource(HikariConfig().apply {
            jdbcUrl = if (readOnly) "jdbc:sqlite:file:$dbFile?mode=ro&uri=true" else "jdbc:sqlite:$dbFile"
            driverClassName = "org.sqlite.JDBC"
            maximumPoolSize = poolSize
            connectionInitSql = """
                PRAGMA journal_mode=WAL;
                PRAGMA foreign_keys=ON;
                PRAGMA busy_timeout=$busyTimeout;
            """.trimIndent()
        })

        reader = Database.connect(hikari(readOnly = true, poolSize = readPoolSize))
        writer = Database.connect(hikari(readOnly = false, poolSize = 1))
    }
}

suspend fun <T> dbRead(block: Transaction.() -> T): T =
    withContext(AppDatabase.readDispatcher) {
        transaction(AppDatabase.reader) { block() }
    }

suspend fun <T> dbWrite(block: Transaction.() -> T): T =
    withContext(AppDatabase.writeDispatcher) {
        transaction(AppDatabase.writer) { block() }
    }
