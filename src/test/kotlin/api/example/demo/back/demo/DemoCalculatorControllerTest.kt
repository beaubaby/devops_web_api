package api.example.demo.back.demo

import io.kotlintest.shouldBe
import io.kotlintest.specs.FreeSpec
import org.junit.jupiter.api.extension.ExtendWith
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.MediaType
import org.springframework.test.context.ActiveProfiles
import org.springframework.test.context.junit.jupiter.SpringExtension
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post

@ExtendWith(SpringExtension::class)
@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("test")
class DemoCalculatorControllerTest(
    @Autowired val mockMvc: MockMvc
) : FreeSpec() {
    init {
        "/demo" - {
            "/" {
                val result = mockMvc.perform(get("/demo"))
                result.andExpect {
                    it.response.contentAsString shouldBe "Hello Calculator"
                }
            }

            "/calculator/plus" {
                val result = mockMvc.perform(
                        post("/demo/calculator/plus")
                                .contentType(MediaType.APPLICATION_JSON)
                                .content("""
                                    {
                                        "a": "7",
                                        "b": "6"
                                    }
                                """.trimIndent())
                )
                result.andExpect {
                    it.response.contentAsString shouldBe "13"
                }
            }
        }
    }
}
