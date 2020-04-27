package com.ka.loaneligibility.backend.repositories

import com.ka.loaneligibility.backend.models.CarModel
import javax.transaction.Transactional
import org.springframework.data.repository.CrudRepository

@Transactional
interface CarRepository : CrudRepository<CarModel, Long>
