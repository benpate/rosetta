package compare

// OperatorGreaterThan represents an "greater than" comparison, when used in Predicates and Criteria
const OperatorGreaterThan = ">"

// OperatorGreaterOrEqual represents an "greater or equal" comparison, when used in Predicates and Criteria
const OperatorGreaterOrEqual = ">="

// OperatorEqual represents an "equals" comparison, when used in Predicates and Criteria
const OperatorEqual = "="

// OperatorNotEqual represents a "not equals" comparison, when used in Predicates and Criteria
const OperatorNotEqual = "!="

// OperatorLessOrEqual represents an "less or equal" comparison, when used in Predicates and Criteria
const OperatorLessOrEqual = "<="

// OperatorLessThan represents a "less than" comparison, when used in Predicates and Criteria
const OperatorLessThan = "<"

// OperatorBeginsWith represents a "begins with" comparison , when used in Predicates and Criteria.  It is only valid for strings and []string.
const OperatorBeginsWith = "BEGINS"

// OperatorEndsWith represents a "ends with" comparison, when used in Predicates and Criteria.  It is only valid for strings and []string.
const OperatorEndsWith = "ENDS"

// OperatorContains represents a "contains" comparison, when used in Predicates and Criteria.  It is valid for the types that Contains accepts.
const OperatorContains = "CONTAINS"

// OperatorContainedBy represents a "contains" comparison in reverse, when used in Predicates and Criteria.  It is valid for the types that Contains accepts.
const OperatorContainedBy = "CONTAINED BY"
