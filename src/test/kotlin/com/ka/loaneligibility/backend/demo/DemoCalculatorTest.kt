package com.ka.loaneligibility.backend.demo

import io.kotlintest.data.forall
import io.kotlintest.shouldBe
import io.kotlintest.specs.FreeSpec
import io.kotlintest.tables.row

class DemoCalculatorTest : FreeSpec({
    "plus" - {
        val firstRandomInt = Math.random().toInt()
        val secondRandomInt = Math.random().toInt()

        val demoCalculator = DemoCalculator()
        forall(
                row(1, 1, 2),
                row(5, 5, 10),
                row(firstRandomInt, secondRandomInt, firstRandomInt + secondRandomInt)
        ) { a, b, result ->
            demoCalculator.plus(a, b) shouldBe result
        }
    }
})
