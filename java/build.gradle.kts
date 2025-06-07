plugins {
	java
	id("org.springframework.boot") version "3.4.6"
	id("io.spring.dependency-management") version "1.1.7"
	id("com.google.protobuf") version "0.9.5"
}

group = "centraluniversity.app"
version = "0.0.1-SNAPSHOT"

java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(17)
	}
}

configurations {
	compileOnly {
		extendsFrom(configurations.annotationProcessor.get())
	}
}

repositories {
	mavenCentral()
}

dependencies {
	implementation("org.springdoc:springdoc-openapi-starter-webmvc-ui:2.8.8")
	implementation("org.springframework.boot:spring-boot-starter-jdbc")
	implementation("org.springframework.boot:spring-boot-starter-web")
	implementation("org.springframework.boot:spring-boot-starter-data-jpa")
	implementation("org.postgresql:postgresql:42.7.2")
	implementation("javax.validation:validation-api:2.0.1.Final")
	
	implementation("com.google.protobuf:protobuf-java:4.30.2")
    implementation("io.grpc:grpc-netty-shaded:1.51.0")
    implementation("io.grpc:grpc-protobuf:1.51.0")
    implementation("io.grpc:grpc-stub:1.51.0")
	implementation("javax.annotation:javax.annotation-api:1.3.2")

	compileOnly("org.projectlombok:lombok")
	annotationProcessor("org.projectlombok:lombok")
	testCompileOnly("org.projectlombok:lombok")
	testAnnotationProcessor("org.projectlombok:lombok")
	
	testImplementation("org.springframework.boot:spring-boot-starter-test")
	testRuntimeOnly("org.junit.platform:junit-platform-launcher")
	testImplementation(platform("org.junit:junit-bom:5.10.0"))
    testImplementation("org.junit.jupiter:junit-jupiter")
}

protobuf {
	protoc {
        artifact = "com.google.protobuf:protoc:3.21.12"  
    }
    plugins {
        create("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:1.51.0"
        }
    }
    generateProtoTasks {
        all().forEach {
            it.plugins {
                create("grpc")
            }
        }
    }
}

tasks.withType<Test> {
	useJUnitPlatform()
}
