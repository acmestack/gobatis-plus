package constants

const SELECT_SQL = "SELECT #{columns} FROM #{tableName} WHERE #{conditions}"

const INSERT_SQL = "INSERT INTO #{tableName} (#{columns}) VALUES (#{columnMapping})"

const UPDATEBYID_SQL = "UPDATE #{tableName} SET #{columnMapping} WHERE #{conditions}"

const DELETEBYID_SQL = "delete from #{tableName} where #{conditions}"
