package constants

const SELECT_SQL = "SELECT #{columns} FROM #{tableName} WHERE #{conditions}"

const INSERT_SQL = "INSERT INTO #{tableName} (#{columns}) VALUES (#{columnMapping})"
const INSERT_BATCH_SQL = "INSERT INTO #{tableName} (#{columns}) VALUES (#{columnMapping}),(#{columnMapping}),(#{columnMapping})"

const DELETEBYID_SQL = "delete from #{tableName} where #{conditions}"
const DELETEBATCHIDS_SQL = "delete from #{tableName} where in #{conditions}"

const UPDATEBYID_SQL = "update #{tableName} set #{filed} where #{conditions}"
