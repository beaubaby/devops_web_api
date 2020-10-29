package api.example.demo.back.repositories

import api.example.demo.back.models.CarModel
import javax.transaction.Transactional
import org.springframework.data.repository.CrudRepository

@Transactional
interface CarRepository : CrudRepository<CarModel, Long>
