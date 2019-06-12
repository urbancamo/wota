include_once("dbconfig.php");
include_once("include/errstr.php");
include_once("include/ucall.php");
function stncall($call) {
  $ucall = strtoupper($call);
  $p = strpos($ucall,"/P");
  if($p > 2) {
    $ucall = substr($ucall,0,$p);
  }
  return $ucall;
}
$self = $_SERVER['PHP_SELF']."?page=".$_REQUEST['page'];
// initialization
global $gCms;
$feu =& $gCms->modules['FrontEndUsers']['object'];
$logged_by = $feu->LoggedInName();
$callused = $logged_by;
global $errcount;
$errcount = 0;
if($_SERVER['REQUEST_METHOD'] == "POST") {
  // load the form fields
  $date = isset($_POST['date1']) ? mysql_real_escape_string($_POST['date1']) : "";
  $year = substr($date,0,4);
  $call_used = strtoupper(trim(mysql_real_escape_string($_POST['call_used'])));
  $summitid = mysql_real_escape_string($_POST['summit']);
  $stncall = stncall(strtoupper(trim(mysql_real_escape_string($_POST['call']))));
  $wkdby = strtoupper(trim(mysql_real_escape_string($_POST['logged_by'])));
  if($stncall=="") {
    $message = errstr("Call missing");
  }
  switch($year) {
    case '2009': break;
    case '2010': break;
    case '2011': break;
    case '2012': break;
    case '2013': break;
    case '2014': break;
    case '2015': break;
    case '2016': break;
    case '2017': break;
    case '2018': break;
    case '2019': break;


    default: $message = errstr("Cannot log contacts for year ".$year);
  }
}
if(($_SERVER['REQUEST_METHOD'] == "POST")&&($errcount==0)) {
  // check all time points
  $wawpoints = 1; $points = 1;
  $sql = "SELECT * FROM `chaser_log` WHERE `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."'";
  $result = mysql_query($sql,$con);
  if($result && (mysql_num_rows($result) > 0)) $wawpoints = 0;
  $sql = "SELECT * FROM `chaser_log` WHERE `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$stncall."'";
  $result = mysql_query($sql,$con);
  if($result && (@mysql_num_rows($result) > 0)) $points = 0;
  // check annual points
  $wawpoints_yr = 1; $points_yr = 1;
  $sql = "SELECT * FROM `chaser_log` WHERE `year` = ".$year." AND `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."'";
  $result = mysql_query($sql,$con);
  if($result && (mysql_num_rows($result) > 0)) $wawpoints_yr = 0;
  $sql = "SELECT * FROM `chaser_log` WHERE `year` = ".$year." AND `wkdby` = '".$wkdby."' AND `wotaid` = '".$summitid."' AND `stncall` = '".$stncall."'";
  $result = mysql_query($sql,$con);
  if($result && (@mysql_num_rows($result) > 0)) $points_yr = 0;
  // check confirmation
  $sql = "SELECT * FROM `activator_log` WHERE `callused` = '".$stncall."' AND `wotaid` = '".$summitid."' AND `ucall` = '".ucall($wkdby)."' AND `date` = '".$date."'";
  $result = mysql_query($sql,$con);
  if($result && (@mysql_num_rows($result) > 0)) {
    $cfmmsg = " confirmed and";
    mysql_query("UPDATE `activator_log` SET `confirmed` = true WHERE `callused` = '".$stncall."' AND `wotaid` = '".$summitid."' AND `ucall` = '".$wkdby."' AND `date` = '".$date."'");
    $cfm = "true";
    // update
  } else {
    $cfmmsg = "";
    $cfm = "false";
  }
  // log contact
  $sql = "INSERT into `chaser_log` (`wkdby`, `ucall`, `wotaid`, `date`, `year`, `stncall`, `points`, `wawpoints`, `points_yr`, `wawpoints_yr`, `confirmed`) VALUES ('".$wkdby."','".$call_used."',".$summitid.",'".$date."',".$year.",'".$stncall."', ".$points.",".$wawpoints.",".$points_yr.",".$wawpoints_yr.",".$cfm.")";
  $result = mysql_query($sql,$con);
?>
<h3>Contact<?=$cfmmsg?> logged</h3>
<table width="50%" border="0" cellspacing="8" cellpadding="8">
<tr><th width="60%" align="center">Points</th><th width="20%" align="center"><?=$year?></th><th width="20%" align="center">All-time</th></tr>
<tr><td>Peak bagger points</td><td align="center"><?=$wawpoints_yr?></td><td align="center"><?=$wawpoints?></td></tr>
<tr><td>Fell runner points</td><td align="center"><?=$points_yr?></td><td align="center"><?=$points?></td></tr>
</table>
<p><a href="<?=$self?>">Log another contact</a></p>
<p><a href="mm_home">Go to your home page</a></p>
<?
} else {
  require_once('classes/tc_calendar.php');
  //instantiate class and set properties
  $myCalendar = new tc_calendar("date1", true);
  $myCalendar->setIcon("images/iconCalendar.gif");
  $myCalendar->setDate(date(d),date(m),date(Y));
  if($errcount > 0) echo $message;
// display the form
?>
<script language="javascript" src="calendar.js"></script>
<form action="<?=$self?>" method="POST">
    <input type="hidden" name="logged_by" value="<?=$feu->LoggedInName()?>">
    <table border="0" cellpadding="2">
        <tr>
            <td>Date of contact:&nbsp;</td>
            <td><? $myCalendar->writeScript(); ?></td>
        </tr>
        <tr>
            <td>Callsign used:&nbsp;</td>
            <td><input type="text" size="8" name="call_used" value="<?=$callused?>"></td>
        </tr>
        <tr>
            <td width="25%">Summit:</td>
            <td width="75%">
<select name="summit">
<?
$result = mysql_query("SELECT `wotaid`,`sotaid`,`name` FROM `summits` WHERE `wotaid` > 0 ORDER BY `name`",$con);
if (mysql_num_rows($result) ) {
  while($row = mysql_fetch_array($result, MYSQL_ASSOC)) {
    $val = $row['wotaid'] + 0;

  //  $wotaid = "LDW-".substr($row['wotaid'],-3,3);

$fellid1 = substr($row['wotaid'],-3,3);
  $fellid = (int)$fellid1;
  if ($fellid <= 214) {$prefix = "LDW-";}
  else {$prefix = "LDO-";
  $fellid = (int)$fellid-214;}
  $ldid = $prefix . str_pad($fellid,3,"0",STR_PAD_LEFT);
 $wotaid = $ldid;
 $sotaid= "";
  if ($row['sotaid'] != '') {
   $sotano = (int)$row['sotaid'];
   $sotaref = "G/LD-" . str_pad($sotano,3,"0",STR_PAD_LEFT);
   $sotaid = "[". $sotaref . "]";
  }
    echo "<option value=\"".$val."\">".$row['name']." (".$wotaid.") ".$sotaid." </option>";
  }
}
?>
</select>
            </td>
        </tr>
        <tr>
            <td>Station worked:</td>
            <td><input type="text" size="8" name="call"></td>
        </tr>
    </table>
    <p><input type="submit" value="Submit">&nbsp;&nbsp;<button type="button" onclick="javascript:window.location='mm_home.html'">Cancel</button></p>
</form>
<?
}
