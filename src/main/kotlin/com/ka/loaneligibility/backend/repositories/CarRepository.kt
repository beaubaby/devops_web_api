package com.ka.loaneligibility.backend.repositories

import com.ka.loaneligibility.backend.models.CarModel
import org.springframework.data.repository.CrudRepository
import javax.transaction.Transactional

@Transactional
interface CarRepository : CrudRepository<CarModel, Long> {
}