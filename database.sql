-- phpMyAdmin SQL Dump
-- version 4.4.10
-- http://www.phpmyadmin.net
--
-- Host: localhost:3306
-- Generation Time: Jan 25, 2017 at 12:31 PM
-- Server version: 5.5.42
-- PHP Version: 5.6.10

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

--
-- Database: `golangblog`
--

-- --------------------------------------------------------

--
-- Table structure for table `categories`
--

CREATE TABLE `categories` (
  `categorie_id` int(255) NOT NULL,
  `title` varchar(150) NOT NULL,
  `description` varchar(160) NOT NULL,
  `content` varchar(5000) NOT NULL,
  `keywords` varchar(3000) NOT NULL,
  `approved` tinyint(4) NOT NULL DEFAULT '0',
  `author` varchar(30) NOT NULL,
  `type` varchar(15) NOT NULL,
  `date` date NOT NULL,
  `parent_id` int(30) DEFAULT '0',
  `trashed` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `categories`
--

INSERT INTO `categories` (`categorie_id`, `title`, `description`, `content`, `keywords`, `approved`, `author`, `type`, `date`, `parent_id`, `trashed`) VALUES
(32, 'Hello', '', '', '', 0, '', '', '0000-00-00', 0, 0);

-- --------------------------------------------------------

--
-- Table structure for table `files`
--

CREATE TABLE `files` (
  `file_id` int(255) NOT NULL,
  `name` varchar(50) NOT NULL,
  `type` varchar(50) NOT NULL,
  `size` varchar(64) NOT NULL,
  `file_name` varchar(100) NOT NULL,
  `thumb_name` varchar(100) NOT NULL,
  `folder_id` int(255) NOT NULL,
  `date` date NOT NULL,
  `secured` tinyint(1) NOT NULL,
  `path` varchar(5000) NOT NULL,
  `thumb_path` varchar(5000) NOT NULL,
  `user_id` int(255) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `folders`
--

CREATE TABLE `folders` (
  `folder_id` int(255) NOT NULL,
  `folder_name` varchar(60) NOT NULL,
  `description` varchar(140) NOT NULL,
  `author` varchar(60) NOT NULL,
  `parent_id` int(255) DEFAULT '0',
  `path` varchar(1000) NOT NULL,
  `date` date NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `folders`
--

INSERT INTO `folders` (`folder_id`, `folder_name`, `description`, `author`, `parent_id`, `path`, `date`) VALUES
(2, 'Jorn', '', '', 0, 'files/Jorn', '0000-00-00');

-- --------------------------------------------------------

--
-- Table structure for table `pages`
--

CREATE TABLE `pages` (
  `page_id` int(255) NOT NULL,
  `title` varchar(150) NOT NULL,
  `description` varchar(160) NOT NULL,
  `content` varchar(5000) NOT NULL,
  `keywords` varchar(3000) NOT NULL,
  `approved` tinyint(4) NOT NULL DEFAULT '0',
  `author` varchar(30) NOT NULL,
  `date` date NOT NULL,
  `parent_id` int(50) NOT NULL,
  `trashed` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=latin1;

-- --------------------------------------------------------

--
-- Table structure for table `posts`
--

CREATE TABLE `posts` (
  `post_id` int(255) NOT NULL,
  `title` varchar(150) NOT NULL,
  `description` varchar(160) NOT NULL,
  `content` varchar(5000) NOT NULL,
  `keywords` varchar(3000) NOT NULL,
  `approved` tinyint(4) NOT NULL DEFAULT '0',
  `author` varchar(30) NOT NULL,
  `date` date NOT NULL,
  `category_id` int(50) NOT NULL,
  `trashed` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB AUTO_INCREMENT=119 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `posts`
--

INSERT INTO `posts` (`post_id`, `title`, `description`, `content`, `keywords`, `approved`, `author`, `date`, `category_id`, `trashed`) VALUES
(118, 'Hello World!', '', '<p>Hello World!</p>', '', 0, '', '0000-00-00', 32, 0);

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `user_id` int(30) NOT NULL,
  `username` varchar(500) NOT NULL,
  `password` varchar(100) NOT NULL,
  `first_name` varchar(500) NOT NULL,
  `last_name` varchar(500) NOT NULL,
  `dob` varchar(500) NOT NULL,
  `email` varchar(500) NOT NULL,
  `function` varchar(500) NOT NULL,
  `rights` varchar(500) NOT NULL,
  `trashed` tinyint(1) NOT NULL DEFAULT '0',
  `approved` tinyint(1) NOT NULL DEFAULT '0'
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`user_id`, `username`, `password`, `first_name`, `last_name`, `dob`, `email`, `function`, `rights`, `trashed`, `approved`) VALUES
(34, 'admin', '$2a$15$oJzuBoIsAygWc1VGNs6Cpula07g8sfsroLMjoIt9tf8tankTKmWHK', 'Jorn', 'Schalkwijk', '', 'jornschalkwijk@gmail.com', 'Master of the Universe', 'Admin', 0, 0);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`categorie_id`);

--
-- Indexes for table `files`
--
ALTER TABLE `files`
  ADD PRIMARY KEY (`file_id`);

--
-- Indexes for table `folders`
--
ALTER TABLE `folders`
  ADD PRIMARY KEY (`folder_id`);

--
-- Indexes for table `pages`
--
ALTER TABLE `pages`
  ADD PRIMARY KEY (`page_id`);

--
-- Indexes for table `posts`
--
ALTER TABLE `posts`
  ADD PRIMARY KEY (`post_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `categories`
--
ALTER TABLE `categories`
  MODIFY `categorie_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=33;
--
-- AUTO_INCREMENT for table `files`
--
ALTER TABLE `files`
  MODIFY `file_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `folders`
--
ALTER TABLE `folders`
  MODIFY `folder_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=3;
--
-- AUTO_INCREMENT for table `pages`
--
ALTER TABLE `pages`
  MODIFY `page_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=19;
--
-- AUTO_INCREMENT for table `posts`
--
ALTER TABLE `posts`
  MODIFY `post_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=119;
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `user_id` int(30) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=35;