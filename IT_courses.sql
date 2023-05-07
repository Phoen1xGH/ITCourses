-- MySQL dump 10.13  Distrib 8.0.30, for Linux (x86_64)
--
-- Host: localhost    Database: IT courses
-- ------------------------------------------------------
-- Server version	8.0.30

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `CourseCost`
--

DROP TABLE IF EXISTS `CourseCost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CourseCost` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `idCategory` int unsigned NOT NULL,
  `Cost` double unsigned NOT NULL,
  PRIMARY KEY (`id`,`idCategory`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `idCategory_UNIQUE` (`idCategory`),
  KEY `idCategory_idx` (`idCategory`),
  CONSTRAINT `idCategory` FOREIGN KEY (`idCategory`) REFERENCES `DevelopmentCategory` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CourseCost`
--

LOCK TABLES `CourseCost` WRITE;
/*!40000 ALTER TABLE `CourseCost` DISABLE KEYS */;
INSERT INTO `CourseCost` VALUES (1,1,25000),(2,2,27000),(3,3,45000),(4,4,40000),(5,5,55000);
/*!40000 ALTER TABLE `CourseCost` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CourseLanguage`
--

DROP TABLE IF EXISTS `CourseLanguage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CourseLanguage` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ProgrammingLanguage` varchar(45) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `idCategory` int unsigned NOT NULL,
  PRIMARY KEY (`id`,`idCategory`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `id_DevCat_idx` (`idCategory`,`ProgrammingLanguage`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CourseLanguage`
--

LOCK TABLES `CourseLanguage` WRITE;
/*!40000 ALTER TABLE `CourseLanguage` DISABLE KEYS */;
INSERT INTO `CourseLanguage` VALUES (1,'PHP',1),(2,'Python',2),(3,'C#',3),(4,'Java',4),(5,'C++',5);
/*!40000 ALTER TABLE `CourseLanguage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `DevelopmentCategory`
--

DROP TABLE IF EXISTS `DevelopmentCategory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `DevelopmentCategory` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `CategoryName` varchar(45) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `DevelopmentCategory`
--

LOCK TABLES `DevelopmentCategory` WRITE;
/*!40000 ALTER TABLE `DevelopmentCategory` DISABLE KEYS */;
INSERT INTO `DevelopmentCategory` VALUES (1,'Backend'),(2,'QA Automation'),(3,'Desktop Development'),(4,'Mobile Development'),(5,'Game Development');
/*!40000 ALTER TABLE `DevelopmentCategory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `users_id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_displayname` varchar(45) DEFAULT NULL,
  `user_name` varchar(45) DEFAULT NULL,
  `user_pass` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`users_id`),
  UNIQUE KEY `users_id_UNIQUE` (`users_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'Ершов Д.И.','Daniel','$2a$04$NXHFE9mO5eTIeay4ur0OcePHb/svEXVBI/XX6uuk9UReZ8.fVHxNm'),(2,'Некто','N','$2a$04$ACbHpRWOUckvoZz/4/jFhOOpqvpYdZKBGvjadr0y6uG10IEeeQmde');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-02-08 19:50:27
