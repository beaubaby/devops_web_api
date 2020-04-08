package io.kotlintest.provided

import io.kotlintest.AbstractProjectConfig
import io.kotlintest.extensions.ProjectLevelExtension
import io.kotlintest.spring.SpringAutowireConstructorExtension

object ProjectConfig : AbstractProjectConfig() {

    override fun extensions(): List<ProjectLevelExtension> = listOf(SpringAutowireConstructorExtension)

    override fun parallelism(): Int = 10
}
