package com.ka.loaneligibility.backend.models

import javax.persistence.Entity
import javax.persistence.GeneratedValue
import javax.persistence.Id

@Entity
data class CarModel(
    @Id
    @GeneratedValue
    var id: Long? = null,
    var brand: String,
    var color: String
)
