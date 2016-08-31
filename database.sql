-- phpMyAdmin SQL Dump
-- version 4.4.10
-- http://www.phpmyadmin.net
--
-- Host: localhost:3306
-- Generation Time: Aug 06, 2016 at 08:36 PM
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
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `categories`
--

INSERT INTO `categories` (`categorie_id`, `title`, `description`, `content`, `keywords`, `approved`, `author`, `type`, `date`, `parent_id`, `trashed`) VALUES
(3, 'Mensen', '', '', '', 1, 'admin', 'product', '0000-00-00', 0, 0),
(4, 'Apen', '', '', '', 1, 'admin', 'product', '0000-00-00', 0, 0),
(7, 'hufter', 'hufters zijn hufterig', '', '', 1, '', '', '0000-00-00', 0, 0),
(8, 'Penkie', 'henkie', '', '', 1, '', '', '0000-00-00', 0, 0),
(9, 'Schaap', '', '', '', 1, '', '', '0000-00-00', 0, 0),
(10, 'Homo', '', '', '', 0, '', '', '0000-00-00', 0, 0);

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
  `album_id` int(255) NOT NULL,
  `date` date NOT NULL,
  `secured` tinyint(1) NOT NULL,
  `path` varchar(5000) NOT NULL,
  `thumb_path` varchar(5000) NOT NULL,
  `user_id` int(255) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `files`
--

INSERT INTO `files` (`file_id`, `name`, `type`, `size`, `file_name`, `thumb_name`, `album_id`, `date`, `secured`, `path`, `thumb_path`, `user_id`) VALUES
(15, 'comments.png', 'png', '0.00', 'August-5-2016-8bdb2a5f-ea93-4c38-7dfc-ea6d28e7833a', '', 0, '0000-00-00', 0, '/file/August-5-2016-8bdb2a5f-ea93-4c38-7dfc-ea6d28e7833a.png', '', 0),
(16, 'Schermafbeelding 2016-08-01 om 15.56.09.png', 'png', '0.42', 'August-5-2016-762069be-9fc0-41a5-473b-3a1605363a76', '', 0, '0000-00-00', 0, '/file/August-5-2016-762069be-9fc0-41a5-473b-3a1605363a76.png', '', 0),
(17, 'Schermafbeelding 2016-08-01 om 22.38.47.png', 'png', '0.34', 'August-6-2016-2ef0dfff-e86a-4557-690e-9aa06cd0333e', '', 0, '0000-00-00', 0, '/file/August-6-2016-2ef0dfff-e86a-4557-690e-9aa06cd0333e.png', '', 0),
(18, 'Schermafbeelding 2016-08-01 om 22.39.01.png', 'png', '0.23', 'August-6-2016-b0b7a1b8-e24a-4a04-624d-93849cc367ba', '', 0, '0000-00-00', 0, '/file/August-6-2016-b0b7a1b8-e24a-4a04-624d-93849cc367ba.png', '', 0),
(19, 'Schermafbeelding 2016-08-01 om 15.56.09.png', 'png', '0.42', 'August-6-2016-220a3375-3cb8-4c73-7168-5f608f5c22d7', '', 0, '0000-00-00', 0, '/file/August-6-2016-220a3375-3cb8-4c73-7168-5f608f5c22d7.png', '', 0),
(20, 'Schermafbeelding 2016-08-01 om 15.56.09.png', 'png', '0.42', 'August-6-2016-c0f29a51-ab8c-4035-5ee2-99cf8c899ea1', '', 0, '0000-00-00', 0, '/file/August-6-2016-c0f29a51-ab8c-4035-5ee2-99cf8c899ea1.png', '', 0),
(21, 'Schermafbeelding 2016-08-02 om 08.54.45.png', 'png', '0.22', 'August-6-2016-12592f5b-1d99-4668-575c-58f99cb98e00', '', 0, '0000-00-00', 0, '/file/August-6-2016-12592f5b-1d99-4668-575c-58f99cb98e00.png', '', 0),
(22, 'Schermafbeelding 2016-08-02 om 08.54.45.png', 'png', '0.21', 'August-6-2016-a04847c0-752e-4113-52e9-0c3d28e45116', '', 0, '0000-00-00', 0, '/file/August-6-2016-a04847c0-752e-4113-52e9-0c3d28e45116.png', '', 0);

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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `folders`
--

INSERT INTO `folders` (`folder_id`, `folder_name`, `description`, `author`, `parent_id`, `path`, `date`) VALUES
(8, 'test', '', '', 0, 'files/test', '0000-00-00'),
(9, 'jorn', '', '', 0, 'files/jorn', '0000-00-00');

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
) ENGINE=InnoDB AUTO_INCREMENT=102 DEFAULT CHARSET=latin1;

--
-- Dumping data for table `posts`
--

INSERT INTO `posts` (`post_id`, `title`, `description`, `content`, `keywords`, `approved`, `author`, `date`, `category_id`, `trashed`) VALUES
(29, 'Test', '', 'jorn', '', 0, 'admin', '2015-12-17', 4, 0),
(30, 'help', '', 'help mij nu', '', 0, 'admin', '2015-12-17', 7, 0),
(31, 'hallo', '', 'wereld', '', 1, '', '0000-00-00', 7, 0),
(32, 'hello', '', 'world', '', 1, '', '0000-00-00', 7, 0),
(33, 'jaapie', '', '<p>is aap</p>', '', 0, '', '0000-00-00', 4, 0),
(96, 'henkie', '', 'is hamster', '', 0, '', '0000-00-00', 7, 0),
(97, 'Jorn', '', '<p>Jorn is een baas</p>', '', 0, '', '0000-00-00', 4, 0),
(98, 'Lotte', '', 'Lekkertje', '', 0, '', '0000-00-00', 7, 0),
(99, 'flololo', 'konijn', '<p>lolo</p>', '', 0, '', '0000-00-00', 9, 0),
(100, 'kankerhomo', '', '<p>ddd</p>', '', 0, '', '0000-00-00', 10, 0),
(101, 'Jornkanker', 'kankerjorn', '<p>fgfg</p>', '', 0, '', '0000-00-00', 9, 0);

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
(33, 'henki', '$2a$10$Xabrtzkh7zQTyxa2JHnYX.78gFcpvoJG8nvCv/wshv3XBMQ1Lk9I6', 'henkie', 'penkie', '', 'poepchinees@gmail.com', 'hamster', '', 0, 0),
(34, 'admin', '$2a$10$3VizguFphwAkxMWjI0VWsuEuxoGvFWpZZZdwSVNARmlgSkTfVXFOW', 'Jorn', 'Schalkwijk', '', 'jornschalkwijk@gmail.com', 'Master of the Universe', 'Admin', 0, 0);

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
  MODIFY `categorie_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=11;
--
-- AUTO_INCREMENT for table `files`
--
ALTER TABLE `files`
  MODIFY `file_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=23;
--
-- AUTO_INCREMENT for table `folders`
--
ALTER TABLE `folders`
  MODIFY `folder_id` int(255) NOT NULL AUTO_INCREMEN;
--
-- AUTO_INCREMENT for table `posts`
--
ALTER TABLE `posts`
  MODIFY `post_id` int(255) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=102;
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `user_id` int(30) NOT NULL AUTO_INCREMENT,AUTO_INCREMENT=35;